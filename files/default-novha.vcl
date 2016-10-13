vcl 4.0;

# Default backend definition. Set this to point to your content server.
backend default {
    .host = "backend";
    .port = "8080";
}

sub vcl_recv {
}

sub vcl_backend_response {
}

sub vcl_deliver {
}
