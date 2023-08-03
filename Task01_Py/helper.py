import random
from datetime import datetime
from typing import Dict, Any

id = 0

def generate_data() -> Dict[str, Any]:
    global id
    id += 1
    data = {
        "id": id,
        "ts": datetime.now().timestamp(),
        "task01": datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")[:-3],
        "payload": {
            "data1": random.randint(1, 100),
            "data2": "TEST DATA with id " + str(id)
        },
        "items": random.choices(["A", "B", "C", "D", "E"], k=3)
    }
    return data