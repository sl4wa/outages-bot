FROM python:3.11-slim

RUN apt-get update && apt-get install -y cron tzdata && rm -rf /var/lib/apt/lists/*
ENV TZ=Europe/Kyiv

WORKDIR /app

COPY . /app

RUN pip install --upgrade pip && pip install -r requirements.txt

RUN printf "CRON_TZ=Europe/Kyiv\n*/10 * * * * root cd /app && /usr/local/bin/python -u /app/notifier.py >> /app/cron.log 2>&1\n" > /etc/cron.d/notifier && \
    chmod 0644 /etc/cron.d/notifier && \
    crontab /etc/cron.d/notifier

COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

CMD ["/bin/sh", "/app/start.sh"]
