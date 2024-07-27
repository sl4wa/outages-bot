import pytz
from datetime import datetime
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters, ConversationHandler
from commands.start import start, street_selection, building_selection
from commands.subscription import show_subscription
from commands.stop import handle_stop
from services.fetch_notify import fetch_and_notify
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.executors.pool import ThreadPoolExecutor
from apscheduler.jobstores.memory import MemoryJobStore
import logging

STREET, BUILDING = range(2)
scheduler = None

# Time constants
STOP_TIME = "00:00"  # Stop time in 24-hour format
START_TIME = "07:00"  # Start time in 24-hour format
UKRAINE_TZ = pytz.timezone('Europe/Kiev')

# Parse time constants
stop_hour, stop_minute = map(int, STOP_TIME.split(':'))
start_hour, start_minute = map(int, START_TIME.split(':'))

def stop_bot(updater):
    updater.stop()
    logging.info("Bot polling stopped.")

def start_bot(updater):
    updater.start_polling()
    logging.info("Bot polling started.")

def disable_fetch_and_notify():
    global scheduler
    job = scheduler.get_job('fetch_and_notify')
    if job:
        scheduler.remove_job('fetch_and_notify')
        logging.info("fetch_and_notify job disabled.")

def enable_fetch_and_notify():
    global scheduler
    job = scheduler.get_job('fetch_and_notify')
    if not job:
        scheduler.add_job(fetch_and_notify, 'interval', minutes=5, id='fetch_and_notify')
        logging.info("fetch_and_notify job enabled.")

def setup_bot(token):
    global scheduler

    jobstores = {'default': MemoryJobStore()}
    executors = {'default': ThreadPoolExecutor(20)}
    job_defaults = {'coalesce': False, 'max_instances': 1}
    scheduler = BackgroundScheduler(jobstores=jobstores, executors=executors, job_defaults=job_defaults, timezone=UKRAINE_TZ)

    # Start the scheduler and add the fetch_and_notify job
    scheduler.start()
    logging.info("Scheduler started.")
    scheduler.add_job(fetch_and_notify, 'interval', minutes=5, id='fetch_and_notify')

    # Run fetch_and_notify on start with timezone awareness
    current_time = datetime.now(UKRAINE_TZ)
    stop_time = UKRAINE_TZ.localize(datetime.combine(current_time.date(), datetime.strptime(STOP_TIME, "%H:%M").time()))
    start_time = UKRAINE_TZ.localize(datetime.combine(current_time.date(), datetime.strptime(START_TIME, "%H:%M").time()))

    if start_time <= current_time < stop_time:
        logging.info("Current time is within restricted hours. fetch_and_notify will not run.")
    else:
        logging.info("Running fetch_and_notify on startup...")
        fetch_and_notify()
        logging.info("fetch_and_notify has been executed on startup.")

    updater = Updater(token, use_context=True)
    dp = updater.dispatcher

    start_conv_handler = ConversationHandler(
        entry_points=[CommandHandler('start', start)],
        states={
            STREET: [MessageHandler(Filters.text & ~Filters.command, street_selection)],
            BUILDING: [MessageHandler(Filters.text & ~Filters.command, building_selection)]
        },
        fallbacks=[],
        allow_reentry=True
    )

    dp.add_handler(start_conv_handler)
    dp.add_handler(CommandHandler('subscription', show_subscription))
    dp.add_handler(CommandHandler('stop', handle_stop))

    # Schedule fetch_and_notify to be disabled and enabled at specified times
    scheduler.add_job(disable_fetch_and_notify, 'cron', hour=stop_hour, minute=stop_minute, id='disable_fetch_and_notify', timezone=UKRAINE_TZ)
    scheduler.add_job(enable_fetch_and_notify, 'cron', hour=start_hour, minute=start_minute, id='enable_fetch_and_notify', timezone=UKRAINE_TZ)

    logging.info("Bot setup completed. Starting polling...")
    updater.start_polling()
    updater.idle()
    logging.info("Bot is now polling.")

def get_scheduler():
    global scheduler
    return scheduler
