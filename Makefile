.PHONY: srv-image
srv-image:
	docker build -t node-crawler-srv:latest .

.PHONY: web-image
web-image:
	docker build -t node-crawler-web:latest ./frontend


UID := $(id -u)
GUID := $(id -g)
export UID
export GID


.PHONY: start
start:
	docker-compose up -d

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
	    mysql -hdb --port=3306 -uncrawl -ptmppswd crawler
