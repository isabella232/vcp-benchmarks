vcl 4.0;
import directors;

backend varnish0 {
    .host = "varnish0";
    .port = "6081";
}

backend varnish1 {
    .host = "varnish1";
    .port = "6081";
}

backend varnish2 {
    .host = "varnish2";
    .port = "6081";
}

backend varnish3 {
    .host = "varnish3";
    .port = "6081";
}


sub vcl_init {
    new vcplb = directors.round_robin();
    vcplb.add_backend(varnish0);
    vcplb.add_backend(varnish1);
    vcplb.add_backend(varnish2);
    vcplb.add_backend(varnish3);
}

sub vcl_recv {
    set req.backend_hint = vcplb.backend();
    return (pass);
}

sub vcl_backend_response {
}

sub vcl_deliver {
}
