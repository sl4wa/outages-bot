import pytz
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters, ConversationHandler
from commands.start import start, street_selection, building_selection
from commands.subscription import show_subscription
from commands.stop import handle_stop
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.executors.pool import ThreadPoolExecutor
from apscheduler.jobstores.memory import MemoryJobStore
import logging

STREET, BUILDING = range(2)
scheduler = None
UKRAINE_TZ = pytz.timezone('Europe/Kiev')

def setup_bot(token):
    global scheduler

    jobstores = {'default': MemoryJobStore()}
    executors = {'default': ThreadPoolExecutor(20)}
    job_defaults = {'coalesce': False, 'max_instances': 1}
    scheduler = BackgroundScheduler(jobstores=jobstores, executors=executors, job_defaults=job_defaults, timezone=UKRAINE_TZ)

    scheduler.start()
    logging.info("Scheduler started.")

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

    logging.info("Bot setup completed. Starting polling...")
    updater.start_polling()
    updater.idle()
    logging.info("Bot is now polling.")

def get_scheduler():
    global scheduler
    return scheduler
