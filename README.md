Nomad and Vault Simulator for testing a CI deployer

This repo includes others work in these directories:

    minica
        from: https://github.com/jsha/minica

    tls-certificates
        from: https://github.com/PrakharSrivastav/tls-certificates

The above directories are great tutorials on PKI certificates creation and use in 'go'.

I have also included PDF's of the authors blogs

-=-=-

The starting point for the code in 'https-nomad-vault-simulations' was copied from:

    tls-certificates\03-https-client\server

I have then adapted it to simulate various Nomad and Vault endpoints for testing a CI deployer.

To run my Nomad/Vault mock, simply cd into: https-nomad-vault-simulations

and run with:

    go run server.go


NOTE:
My code simply uses the authors original certificates that have been pulled from github and these
appear to be OK for the tests i have done using this code with a CI deployer ... in so far as there
have been no complaints about the certificates being expired.
I say this as the 'minica' go application says it generates certificates with an expiry of 2 years
and 30 days to satisfy some apple ios requirement.
If problems are seen, then for any future testing, new certificates, etc will need to be generated,
requiring reading the included readme's in the other included directories.
