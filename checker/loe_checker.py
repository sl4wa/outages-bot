import requests

from .checker_interface import CheckerInterface

API_URL = "https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693"

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
    "Accept": "application/json, text/plain, */*",
    "Connection": "keep-alive",
    "Accept-Language": "en-US,en;q=0.9",
}


class LOEChecker(CheckerInterface):
    """Concrete implementation of CheckerInterface for LOE outages."""

    def get_outages(self):
        """
        Fetch outage data from the LOE API and return cleaned data.
        """
        response = requests.get(API_URL, headers=HEADERS)
        if response.status_code == 200:
            data = response.json()
            outages = data.get("hydra:member", [])

            # Clean and structure the outages
            cleaned_outages = [
                {
                    "dateEvent": outage.get("dateEvent"),
                    "datePlanIn": outage.get("datePlanIn"),
                    "city": outage.get("city", {}),
                    "street": outage.get("street", {}),
                    "buildingNames": outage.get("buildingNames"),
                    "koment": outage.get("koment", ""),
                }
                for outage in outages
            ]
            return cleaned_outages
        else:
            raise ValueError(
                f"Failed to fetch data: HTTP {response.status_code}")
