from abc import ABC, abstractmethod

class CheckerInterface(ABC):
    """Interface for fetching and processing outages."""
    
    @abstractmethod
    def get_outages(self):
        """
        Fetch outage data and return a cleaned list of outages.
        """
        pass
