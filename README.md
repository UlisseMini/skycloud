# skycloud

CLI access to [skycloud](https://skycloud.pro) aka free file hosting :)

Basically skycloud give you a ton of free storage but you can't upload till you upgrade to a paid plan, they only enforce this client side so here you go, now you can upload with a free account and get free storage.


# Bugs

## certificate signed by unknown authority

You need to import the letsencrypt cert from [here](https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem.txt), for example on archlinux
```sh
curl https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem.txt -o letsencrypt.cert
sudo trust anchor --store letsencrypt.cert
```

## Will give an invalid link if you uploaded an image or video

In the link change `Files` to `Video`, `Images` or `Audio` respectively.
