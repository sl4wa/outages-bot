import json

def minimize_json(file_path):
    with open(file_path, 'r', encoding='utf-8') as infile:
        data = json.load(infile)
    
    if "hydra:member" in data:
        minimized_members = [
            {"id": member["id"], "name": member["name"]}
            for member in data["hydra:member"]
        ]
        minimized_data = {"hydra:member": minimized_members}
    else:
        minimized_data = {}

    with open(file_path, 'w', encoding='utf-8') as outfile:
        json.dump(minimized_data, outfile, ensure_ascii=False, indent=4)

if __name__ == "__main__":
    file_path = "streets.json"
    minimize_json(file_path)
