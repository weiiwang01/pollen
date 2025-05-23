### ABSTRACT ###

Pollen is a scalable, high performance, free software (AGPL) web server
that provides small strings of entropy over TLS-encrypted HTTPS or clear
text HTTP connections.  You might think of it as 'Entropy-as-a-Service'.

Pollinate is a free software (GPLv3) script that retrieves entropy from
one or more Pollen servers and seeds the local Pseudo Random Number
Generator (PRNG).  You might think of this as 'PRNG-seeding via
Entropy-as-a-Service'.

Please understand: Neither Pollen nor Pollinate increase the amount of
entropy available on the system!  Rather, Pollinate adequately and
securely seeds the PRNG in cloud virtual machines through communications
with a Pollen server.


### DESCRIPTION ###

The Linux kernel provides two special character devices interfaces to
high quality entropy: /dev/random and /dev/urandom.  Both are pseudo
random number generators (PRNGs), but the former conservatively
guarantees quality entropy, and userspace processes reading from it will
block until sufficient bits are available to fulfill the request.  The
latter, /dev/urandom, provides a non-blocking, limitless stream of
pseudo random numbers.

The manpage random(4) has a far more complete description of /dev/random
and /dev/urandom.  (See [man random](http://manpg.es/random.4))

For most practical purposes /dev/urandom is a perfectly adequate source
of entropy, as long as it is seeded properly at each boot.

Linux distributions (including Debian and Ubuntu) typically carry over a
random seed from one boot to another, typically in an init script such
as /etc/init.d/urandom.  In Ubuntu, that init script does the following:

    ...
    SAVEDFILE=/var/lib/urandom/random-seed
    POOLBYTES=512
    dd if=/dev/urandom of=$SAVEDFILE bs=$POOLBYTES count=1 >/dev/null 2>&1
    ...

There is, you may notice, a bootstrapping problem:  How does one seed
/dev/urandom on a system's very first boot?

On laptops, desktops, tablets, phones, and other physical systems, input
devices such as a keyboard, mouse, touch screen, or microphone can
provide sufficient entropy for seeding the PRNG through the kernel's
collection of timers and interrupts.

However, virtual machines typically have no access to real hardware and
few, if any, sufficient entropy sources.  Several real attacks have been
demonstrated recently against SSH and SSL via certificates generated
with poor entropy, such as:
[Mining Your Ps and Qs: Detection of Widespread Weak Keys in Network Devices](https://factorable.net/weakkeys12.extended.pdf)

The cryptographic security of virtual machines and cloud instances can
be significantly improved by fetching a sufficient amount of entropy at
first boot (and periodically thereafter) to seed the local PRNG with
external sources of entropy.


### IMPLEMENTATION ###

Pollen is a fast and efficient web service implemented in Golang.  It
provides small random strings to its clients over network connections.
Pollen utilizes TLS (SSL) to ensure privacy, security, and
non-repudiation of connections among its clients.

Pollinate is a client utility implemented in Shell, which wraps curl(1)
and communicates securely with one or more Pollen servers.

The default protocol for all connections is HTTPS, however HTTP is
available for debug, testing, and other purposes.  To ensure the privacy
and security of connections the Pollen server should ideally have a
CA-signed certificate, or pre-arrange the distribution of certificates
to its clients.

An entropy request should optimally contain a POST argument:

  - challenge: a randomly generated checksum to ensure unique
    communication with the Pollen server.

The challenge POST argument is a hex-encoded sha512sum(1) value, which
is 128 ASCII characters of [0-9a-f].  Regardless of the value of the
'challenge', the Pollen server will treat the input as a string and
calculate the sha512sum.  This ensures that any malicious input from a
deviant client is whitened to a simple hash before the server operates
upon it.

The server then responds to the client with the sha512 checksum of the
client's challenge on the first line, and the second line of the
response will contain a sha512sum of 64 bytes of entropy.  This second
line is what the client can use as a random seed.

The client verifies the challenge/response, which is intended to help
ensure that this communication between the client and server is a custom
response, and that the server actually "did some work", and thus
affected the entropy state on the server.

The client uses a special option to curl(1) that details all of the
communication to the server, and includes high resolution, local
timestamps.  This information, which is not easily detectable or
reproducible by an attacker (or the Pollen server administrator), is
combined with the server's responses, and written into the Linux PRNG,
/dev/urandom, which is folded into the local system entropy.


### POLLEN AND POLLINATE IN UBUNTU ###

Canonical provides a Pollen server as a service to the Ubuntu community
at [https://entropy.ubuntu.com](https://entropy.ubuntu.com).  Beginning
with Ubuntu 14.04, Ubuntu cloud images include the Pollinate client,
which will try (for up to 3 seconds at first boot) to seed the PRNG with
input from [https://entropy.ubuntu.com](https://entropy.ubuntu.com).
This service is highly available via multiple physical servers deployed
in a cluster using Juju service orchestration.  Each of these Pollen
servers have at least two hardware random number generators, ensuring
high quality entropy as a service, and diversified against hardware
failure.  Moreover, a busy Pollen server, handling many
challenge/response calculations and serving numerous concurrent
connections, will have a computationally complex and impossible to
reproduce entropy state.

Ubuntu cloud users are welcome to add other Pollen servers to their
pool, or just run their own internally, behind their own firewall.
Simply edit the configuration file in /etc/default/pollinate.  Ubuntu
users and other distributions are certainly welcome to install and run
their own Pollen server, with 'sudo apt-get install pollen' or 'bzr
branch lp:pollen' and compile from source.

Be safe, and secure out there!
:-Dustin


### METRICS ###

The metrics subsystem is a new feature in Pollen that you can initialize
and activate using the -metrics-port command-line argument. Once
enabled, the Pollen server will provide standard Prometheus metrics via
the /metrics endpoint at the specified port.

The set of metrics pollen provided is outlined below:

| Metric Name                                       | Metric Type | Metric Description
| ------------------------------------------------- | ----------- | ------------------
| pollen_http_requests_total                        | Counter     | The total number of requests
| pollen_http_responses_codes                       | Counter     | Total responses sent to clients by code
| pollen_http_response_seconds                      | Histogram   | Response time by code
| pollen_system_entropy                             | Gauge       | System available entropy (entropy_avail)
| pollen_response_entropy_per_byte                  | Histogram   | Entropy per bit of the random data in response
| pollen_response_entropy_arithmetic_mean_deviation | Histogram   | Arithmetic mean deviation of the random data in response

Notes:

  - pollen_system_entropy: This metric may not be very useful for
    systems running newer kernels, due to the implementation of the new
    CRNG-based kernel random device.
  - pollen_response_entropy_per_byte and
    pollen_response_entropy_arithmetic_mean_deviation: these two metrics
    might not be very significant as the sample size for each measurement
    is limited to a single response, which is 64 bytes by default.
  - Despite these factors, all three metrics mentioned above continue to
    be observed to be inline with the current pollen logging and
    pollinate's check_pollen script.
