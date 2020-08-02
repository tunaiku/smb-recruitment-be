mkdir database

echo CREATE DATABASE \"mobile-banking-service\" > database/init.sql

cat << 'EOF' >> database/Dockerfile
FROM postgres

COPY init.sql /docker-entrypoint-initdb.d/
EOF

docker build --tag postgres-smb ./database/

rm -rf database

docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name=postgres-smb -d postgres-smb
