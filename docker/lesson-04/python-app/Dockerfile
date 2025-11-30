FROM python:3.12.12-slim-bookworm

ARG WORKDIR=/app
ARG UV_VERSION=0.9.9

RUN pip install uv==${UV_VERSION}

WORKDIR ${WORKDIR}

COPY pyproject.toml uv.lock ${WORKDIR}/

RUN uv sync

COPY . . 

EXPOSE 8080

CMD ["python", "main.py"]
