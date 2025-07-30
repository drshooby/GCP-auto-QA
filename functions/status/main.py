def status_get(request):
    status = request.args.get("status")
    if status:
        return f"Execution status: {status}", 200

    request_json = request.get_json(silent=True)
    if request_json and "status" in request_json:
        return f"Execution status: {request_json["status"]}", 200

    return "Missing 'status' in request", 400
