# Crawlzilla Crawler Telegram Bot

## Description


## Setup And Installation

## How To Run

## Testing

## Project Structure
```mermaid
graph TD;
db1[(Postgres DB)]
redis[(Redis Cache)]

ServiceCrawler-->db1
ServiceCrawler-->Crawler
ServiceCrawler-->Logstash
ServiceCrawler-->redis

ServiceBot-->db1
ServiceBot-->Operation
ServiceBot-->Logstash
ServiceBot-->redis

Elasticsearch-->Kibana
Kibana-->Logstash

Telegram_Bot-->ServiceCrawler
Telegram_Bot-->ServiceBot
Telegram_Bot-->redis
Telegram_Bot-->Commands
Telegram_Bot-->Handlers
Telegram_Bot-->Scenarios
Telegram_Bot-->Logstash

Telegram_Bot_Api<-->Telegram_Bot
```
