ARG UV_VERSION=0.9.9

FROM ghcr.io/astral-sh/uv:${UV_VERSION} AS uv

FROM python:3.12.12-slim-bookworm

COPY --from=uv /uv /uvx /bin/

ARG WORKDIR=/app

WORKDIR ${WORKDIR}

COPY pyproject.toml uv.lock ./

RUN uv sync --frozen --no-dev

COPY . . 

EXPOSE 8080

CMD ["uv", "run", "python", "main.py"]
