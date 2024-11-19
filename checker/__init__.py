from .checker_interface import CheckerInterface
from .loe_checker import LOEChecker

# Global instance of the checker implementation
checker: CheckerInterface = LOEChecker()
