FROM logstash:8.15.3

COPY ./logger/logstash.conf /etc/logstash/conf.d/

CMD logstash -f /etc/logstash/conf.d/logstash.conf