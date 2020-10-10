# SMS service with RabbitMQ

Install RabbitMQ and start broker

```bash
    $ sudo service rabbitmq-server start
    //Service will be started
```

Clone this repo:

```bash
    $ git clone https://github.com/swallow1509/send-sms-via-rabbitmq
    //Go to the directory
```

Run producer/produce.go and consumer/consume.go in two separate terminals

Terminal 1

```bash
    send-sms-via-rabbitmq/producer> go run produce.go Hello, Wolrd!
    //Any message is welcome instead of "Hello, World!"
```

Terminal 2

```bash
    send-sms-via-rabbitmq/consumer> go run consume.go
    //[*] Waiting for messages. To exit press CTRL+C
    //Above message will be displayed untill you pres CTRL+C
```
