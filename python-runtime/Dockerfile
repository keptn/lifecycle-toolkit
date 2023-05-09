FROM python:3.9

# RUN apt-get update -y && \
#     apt-get install git -y

RUN pip install -q --disable-pip-version-check pyyaml GitPython requests

COPY run_script.sh /app/
WORKDIR /app

USER 1000:1000

CMD [ "sh", "-c", "/app/run_script.sh" ]
