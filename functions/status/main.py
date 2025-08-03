import os
import requests

def status_get(request):
    request_json = request.get_json(silent=True)
    
    if request_json and "results" in request_json:
        slack_webhook_url = os.getenv("SLACK_WEBHOOK_URL")
        
        if not slack_webhook_url:
            return "SLACK_WEBHOOK_URL not set", 500
        
        results = request_json["results"]
        
        passed_count = sum(1 for result in results if result.get("passed", False))
        total_count = len(results)
        
        if passed_count == total_count:
            text = f"✅ All {total_count} tests passed!"
        else:
            text = f"❌ {passed_count}/{total_count} tests passed"
        
        details = []
        for result in results:
            status_icon = "✅" if result.get("passed", False) else "❌"
            test_name = result.get("name", "Unknown test")
            details.append(f"{status_icon} {test_name}")
            
            if not result.get("passed", False) and result.get("errors"):
                for error in result["errors"]:
                    details.append(f"    • {error}")
        
        full_message = text + "\n\n" + "\n".join(details)
        
        slack_payload = {"text": full_message}
        
        try:
            response = requests.post(slack_webhook_url, json=slack_payload)
            if response.ok:
                return "Test report sent", 200
            else:
                return f"Slack failed: {response.text}", 502
        except requests.exceptions.RequestException as e:
            return f"Request failed: {e}", 502

    return "Missing 'results' in request", 400