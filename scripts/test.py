import json

import requests

NUMBER_OF_MESSAGES = 2000


class SMS:

    def __init__(self, sender: str, receiver: str, message: str):
        self.sender = sender
        self.receiver = receiver
        self.message = message

    @staticmethod
    def new_from(i: int):
        return SMS('User1', 'User2', 'This is sms number ' + str(i))

    def to_json(self):
        return json.dumps(self, default=lambda o: o.__dict__)


def main():
    for i in range(NUMBER_OF_MESSAGES):
        sms = SMS.new_from(i)
        response = requests.post('http://localhost:8082/api/sms',
                      data=sms.to_json())
        print(response.json())


if __name__ == "__main__":
    main()
