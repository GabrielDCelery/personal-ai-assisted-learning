FROM python:3.14-slim-bookworm

ARG APP_VERSION=1.0.0

ENV ENVIRONMENT=production
ENV PORT=8080

WORKDIR /app
COPY app.py .

LABEL version=${APP_VERSION}

EXPOSE ${PORT}

CMD ["python","-u", "app.py"]

