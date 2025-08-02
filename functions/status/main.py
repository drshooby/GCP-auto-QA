import os
import requests

def status_get(request):
    request_json = request.get_json(silent=True)
    
    if request_json and "results" in request_json:
        slack_webhook_url = os.getenv("SLACK_WEBHOOK_URL")
        
        if not slack_webhook_url:
            return "SLACK_WEBHOOK_URL not set", 500
        
        try:
            response = requests.post(slack_webhook_url, json=request_json)
            if response.ok:
                return "Test report received", 200
            else:
                return f"Slack webhook failed: {response.text}", 502
        except requests.exceptions.RequestException as e:
            return f"Slack request exception: {e}", 502

    return "Missing 'results' in request", 400
