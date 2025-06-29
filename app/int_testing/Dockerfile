FROM postgres:latest

ARG POSTGRES_PASSWORD

ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Copy the fully-initialized Postgres data directory into the image
COPY ../temp/temp_pgdata /var/lib/postgresql/data

# Ensure correct ownership
RUN chown -R postgres:postgres /var/lib/postgresql/data
