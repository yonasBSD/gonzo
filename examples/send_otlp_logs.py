#!/usr/bin/env python3
"""
Example script to send OTLP logs to Gonzo via gRPC or HTTP

Requirements for gRPC:
    pip install opentelemetry-api opentelemetry-sdk opentelemetry-exporter-otlp-proto-grpc

Requirements for HTTP:
    pip install opentelemetry-api opentelemetry-sdk opentelemetry-exporter-otlp-proto-http

Usage:
    1. Start Gonzo with OTLP enabled:
       gonzo --otlp-enabled
       # Or with custom ports:
       gonzo --otlp-enabled --otlp-grpc-port=4317 --otlp-http-port=4318
    
    2. Run this script:
       python send_otlp_logs.py           # Uses gRPC by default
       python send_otlp_logs.py --http    # Use HTTP protocol
"""

import logging
import time
import random
import sys
import argparse

try:
    from opentelemetry._logs import set_logger_provider
    from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
    from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
    from opentelemetry.sdk.resources import Resource
except ImportError:
    print("Error: OpenTelemetry SDK not installed.")
    print("Please install required packages:")
    print("  pip install opentelemetry-api opentelemetry-sdk")
    sys.exit(1)

def get_exporter(use_http=False):
    """Get the appropriate OTLP exporter based on protocol"""
    if use_http:
        try:
            from opentelemetry.exporter.otlp.proto.http._log_exporter import OTLPLogExporter
            print("Using HTTP protocol on port 4318")
            return OTLPLogExporter(
                endpoint="http://localhost:4318/v1/logs",
            )
        except ImportError:
            print("Error: HTTP exporter not installed.")
            print("Please install: pip install opentelemetry-exporter-otlp-proto-http")
            sys.exit(1)
    else:
        try:
            from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
            print("Using gRPC protocol on port 4317")
            return OTLPLogExporter(
                endpoint="localhost:4317",
                insecure=True  # For local testing without TLS
            )
        except ImportError:
            print("Error: gRPC exporter not installed.")
            print("Please install: pip install opentelemetry-exporter-otlp-proto-grpc")
            sys.exit(1)

# Sample log messages with different severities
# Include severity in message text as fallback for parsing
log_messages = [
    ("INFO", "[INFO] Application started successfully"),
    ("INFO", "[INFO] Connecting to database"),
    ("DEBUG", "[DEBUG] Database connection pool initialized"),
    ("INFO", "[INFO] Processing user request"),
    ("WARNING", "[WARN] High memory usage detected"),
    ("ERROR", "[ERROR] Failed to connect to external API"),
    ("INFO", "[INFO] Retrying connection..."),
    ("INFO", "[INFO] Successfully processed request"),
    ("DEBUG", "[DEBUG] Cache hit for key: user_123"),
    ("INFO", "[INFO] Request completed in 245ms"),
    ("TRACE", "[TRACE] Detailed trace information"),
    ("FATAL", "[FATAL] Critical system failure"),
]

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Send OTLP logs to Gonzo')
    parser.add_argument('--http', action='store_true', help='Use HTTP protocol instead of gRPC')
    args = parser.parse_args()
    
    # Configure OpenTelemetry resource
    resource = Resource.create({
        "service.name": "example-python-app",
        "service.version": "1.0.0",
        "environment": "development",
        "host.name": "localhost"
    })
    
    # Get the appropriate exporter
    otlp_exporter = get_exporter(use_http=args.http)
    
    # Set up the logger provider
    logger_provider = LoggerProvider(resource=resource)
    set_logger_provider(logger_provider)
    
    # Add batch processor with OTLP exporter
    logger_provider.add_log_record_processor(
        BatchLogRecordProcessor(otlp_exporter)
    )
    
    # Configure Python logging to use OpenTelemetry
    handler = LoggingHandler(level=logging.NOTSET, logger_provider=logger_provider)
    logging.getLogger().addHandler(handler)
    logging.getLogger().setLevel(logging.INFO)
    
    # Get a logger
    logger = logging.getLogger(__name__)
    
    protocol = "HTTP" if args.http else "gRPC"
    port = 4318 if args.http else 4317
    print(f"Sending OTLP logs to Gonzo via {protocol} on localhost:{port}...")
    print("Make sure Gonzo is running with: gonzo --otlp-enabled")
    print("-" * 60)
    
    try:
        # Send logs with random attributes
        for i in range(20):
            level, message = random.choice(log_messages)
            
            # Add structured attributes
            extra = {
                "request_id": f"req_{i:04d}",
                "user_id": f"user_{random.randint(100, 999)}",
                "latency_ms": random.randint(10, 500),
                "endpoint": random.choice(["/api/users", "/api/products", "/api/orders"]),
                "method": random.choice(["GET", "POST", "PUT", "DELETE"]),
                "status_code": random.choice([200, 201, 400, 404, 500]),
            }
            
            # Log with appropriate level
            if level == "TRACE":
                # Python logging doesn't have TRACE, use DEBUG with lower level
                logger.log(5, message, extra=extra)  # TRACE level
            elif level == "DEBUG":
                logger.debug(message, extra=extra)
            elif level == "INFO":
                logger.info(message, extra=extra)
            elif level == "WARNING":
                logger.warning(message, extra=extra)
            elif level == "ERROR":
                logger.error(message, extra=extra)
            elif level == "FATAL":
                logger.critical(message, extra=extra)
            
            print(f"Sent: [{level}] {message}")
            
            # Small delay between logs
            time.sleep(0.5)
    
    except KeyboardInterrupt:
        print("\nStopping...")
    
    finally:
        # Ensure all logs are flushed
        logger_provider.shutdown()
        print("\nDone sending logs!")

if __name__ == "__main__":
    main()