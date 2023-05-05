FROM python:3.9

COPY run_script.sh /app/
WORKDIR /app

RUN chmod +x /app/run_script.sh

CMD [ "sh", "-c", "/app/run_script.sh" ]