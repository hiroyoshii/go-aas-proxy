docker compose -f e2e/docker-compose.yaml up -d --wait

docker compose -f e2e/docker-compose.yaml exec -T postgres_aas psql -U postgres aas -f /app/aas_schema.sql
docker compose -f e2e/docker-compose.yaml exec -T postgres_aas psql -U postgres aas -f /app/sql/aas_insert.sql

docker compose -f e2e/docker-compose.yaml exec -T postgres_submodel1 psql -U postgres submodel1 -f /app/sql/submodel1_ddl.sql
docker compose -f e2e/docker-compose.yaml exec -T postgres_submodel1 psql -U postgres submodel1 -f /app/sql/submodel1_insert.sql

docker compose -f e2e/docker-compose.yaml exec -T postgres_submodel2 psql -U postgres submodel2 -f /app/sql/submodel2_ddl.sql
docker compose -f e2e/docker-compose.yaml exec -T postgres_submodel2 psql -U postgres submodel2 -f /app/sql/submodel2_insert.sql

docker compose -f e2e/docker-compose.yaml logs
docker compose -f e2e/docker-compose.yaml ps