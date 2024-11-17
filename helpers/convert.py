import os
import json

# Define the path for the existing subscriptions file
old_subscriptions_file = "subscriptions.json"

# Define the directory where new subscription files will be stored
new_data_directory = "users/data"
os.makedirs(new_data_directory, exist_ok=True)

def convert_subscriptions():
    """Convert JSON subscriptions to key-value files."""
    # Load the existing subscriptions from JSON
    try:
        with open(old_subscriptions_file, "r", encoding="utf-8") as file:
            subscriptions = json.load(file)
    except FileNotFoundError:
        print(f"Error: '{old_subscriptions_file}' not found.")
        return

    # Convert each subscription to a separate file
    for chat_id, subscription in subscriptions.items():
        file_path = os.path.join(new_data_directory, f"{chat_id}.txt")
        with open(file_path, "w", encoding="utf-8") as file:
            # Write the subscription in key-value format
            file.write(f"street_id: {subscription.get('street_id', '')}\n")
            file.write(f"street_name: {subscription.get('street_name', '')}\n")
            file.write(f"building: {subscription.get('building', '')}\n")
            file.write(f"start_date: \n")  # Empty fields by default
            file.write(f"end_date: \n")
            file.write(f"comment: \n")
        print(f"Converted subscription for chat_id {chat_id}.")

if __name__ == "__main__":
    convert_subscriptions()
