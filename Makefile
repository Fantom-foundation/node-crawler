.PHONY: start
start:
	docker-compose up --build -d

.PHONY: stop
stop:
	docker-compose stop

.PHONY: clean
clean:
	docker-compose down --volumes

.PHONY: sql
sql:
	docker-compose exec -ti \
	    db \
	    psql -h localhost -p 5432 crawler

