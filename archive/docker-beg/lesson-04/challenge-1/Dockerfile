FROM python:3.12.12-slim-bookworm

COPY --from=ghcr.io/astral-sh/uv:0.9.9 /uv /uvx /bin/

ENV PORT=8080

RUN useradd -m appuser

WORKDIR /app

COPY .python-version pyproject.toml uv.lock ./

RUN uv sync --no-dev --frozen

COPY . .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE ${PORT}

CMD ["uv", "run", "main.py"]


