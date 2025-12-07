FROM python:3.12.12-slim-bookworm

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN useradd -m -u 1000 appuser

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

RUN chown -R appuser:appuser /app

USER appuser

CMD ["python", "app.py"]

# have a better version of this
# FROM ubuntu:22.04
# RUN apt-get update
# RUN apt-get install -y python3
# RUN apt-get install -y python3-pip
# COPY . /app
# WORKDIR /app
# RUN pip3 install -r requirements.txt
# CMD ["python3", "app.py"]
