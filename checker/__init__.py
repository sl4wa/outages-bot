from .loe_checker import LOEChecker
from .checker_interface import CheckerInterface

# Global instance of the checker implementation
checker: CheckerInterface = LOEChecker()