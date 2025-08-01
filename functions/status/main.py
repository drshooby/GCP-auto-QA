def status_get(request):
    request_json = request.get_json(silent=True)
    if request_json and "results" in request_json:
        return "Test report received", 200

    return "Missing 'results' in request", 400