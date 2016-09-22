vcl 4.0;
include "vha.vcl";

# Default backend definition. Set this to point to your content server.
backend default {
    .host = "127.0.0.1";
    .port = "8080";
}

sub vcl_recv {
    call vha_backend_selection;
}

sub vcl_backend_response {
    call vha_clean_headers;
}

sub vcl_deliver {
}
