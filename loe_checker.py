import json
import logging

import requests

from loe_notifier import loe_notifier

API_URL = "https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693"
LOG_FILE = "loe_checker.log"

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
    handlers=[logging.FileHandler(LOG_FILE), logging.StreamHandler()],
)

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
    "Accept": "application/json, text/plain, */*",
    "Connection": "keep-alive",
    "Accept-Language": "en-US,en;q=0.9",
}


def loe_checker():
    logging.info("Attempting to fetch outage data from API...")
    response = requests.get(API_URL, headers=HEADERS)

    if response.status_code == 200:
        try:
            data = response.json()
            outages = data.get("hydra:member", [])

            cleaned_outages = []
            for outage in outages:
                cleaned_outage = {
                    "dateEvent": outage.get("dateEvent"),
                    "datePlanIn": outage.get("datePlanIn"),
                    "city": outage.get("city", {}),
                    "street": outage.get("street", {}),
                    "buildingNames": outage.get("buildingNames"),
                    "koment": outage.get("koment"),
                }
                cleaned_outages.append(cleaned_outage)

            with open("loe_data.json", "w", encoding="utf-8") as f:
                json.dump(cleaned_outages, f, ensure_ascii=False, indent=4)
            logging.info("Outage data saved to loe_data.json")
        except json.JSONDecodeError as e:
            logging.error(f"Error decoding JSON response: {e}")
        except Exception as e:
            logging.error(f"Unexpected error: {e}")

        loe_notifier()
    else:
        logging.error(f"Failed to fetch data: HTTP {response.status_code}")


if __name__ == "__main__":
    loe_checker()
