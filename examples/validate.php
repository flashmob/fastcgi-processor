<?php
# Recipient address validation example
# This does not need to be in a web root, it will be executed directly by php-fpm.
# Simply point the fcgi_script_filename_validate config option to the path of this file
# (making sure that your fcgi has permissions)

// process the request
$rcpt_to = (isset($_GET['rcpt_to'])) ? $_GET['rcpt_to'] : "";
file_put_contents(__DIR__ . "/validate", $rcpt_to, FILE_APPEND);

// Respond with PASSED if everything went OK
echo "PASSED";