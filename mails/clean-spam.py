from pathlib import Path
import re
from sys import argv
import json


domains_blocklist_path = Path(__file__).parent / "spamdomains.json"
domains_blocklist = json.loads(domains_blocklist_path.read_text())

def domains_of_email(email: str):
    domains = set()
    for d in re.finditer(r'https?://([\w.-]+)', email):
        domains.add(d.group(1))
    return domains

def is_spam(email: str) -> bool:
    return len(email.strip()) == 0 or len(domains_blocklist & domains_of_email(email)) > 0

if len(argv) >= 2 and argv[1] == "domains":
    old_blocklist_size = len(domains_blocklist)
    for f in Path(__file__).parent.glob('*.txt'):
        domains_blocklist |= domains_of_email(f.read_text())
    added_count = len(domains_blocklist) - old_blocklist_size
    domains_blocklist_path.write_text(json.dumps(list(domains_blocklist)))
    if added_count != 0:
        print(f"Added {added_count} domains to blocklist")
else:
    trashed_count = 0
    for mailfile in Path(__file__).parent.glob("*.txt"):
        if is_spam(mailfile.read_text()):
            mailfile.rename(mailfile.parent / ".trash" / mailfile.name)
            print(f"Trashed {mailfile.name}")
            trashed_count += 1
    print(f"\nTrashed {trashed_count} mails")
    print(f"Use `{Path(__file__).name} domains` to update spam domains blocklist (by adding all link domains from mails you currently have)")


