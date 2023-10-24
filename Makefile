.PHONY: srv-image
srv-image:
	docker build -t node-crawler-srv:latest .
	docker image prune -f
	docker builder prune -f

.PHONY: web-image
web-image:
	docker build -t node-crawler-web:latest ./frontend
	docker image prune -f
	docker builder prune -f


UID := $(id -u)
GUID := $(id -g)
export UID
export GID

data/mainnet-109331-no-history.g:
	wget -O data/mainnet-109331-no-history.g https://download.fantom.network/mainnet-109331-no-history.g

.PHONY: start
start: data/mainnet-109331-no-history.g
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
