from kafka import KafkaConsumer
import yaml
import json

class Reader:
    def __init__(self, configfile):
        self.config = {}
        with open(configfile) as file:
            self.config = yaml.load(file, Loader=yaml.FullLoader)
        self.kafkaHosts = self.getKafkaHosts()
        self.kafkaTopic = self.config['kafka']['topic']
        self.consumer = KafkaConsumer(bootstrap_servers=self.kafkaHosts, auto_offset_reset='earliest',  fetch_min_bytes=1, fetch_max_wait_ms=500)
        self.consumer.subscribe([self.kafkaTopic])

    # Must already have obtained config file before calling
    def getKafkaHosts(self):
        hostReturn = ""
        for host in self.config['kafka']['hosts']:
            hostReturn += host + ","

        hostReturn = hostReturn[:len(hostReturn)-1]
        return hostReturn
        
    def getConsumer(self):
        return self.consumer
    
