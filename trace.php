<?php
// Debug page to inspect incoming request and headers
header('Content-Type: text/plain; charset=utf-8');
echo "TRACE DEBUG\n";
echo "Request URI: " . ($_SERVER['REQUEST_URI'] ?? '') . "\n";
echo "Host: " . ($_SERVER['HTTP_HOST'] ?? '') . "\n";
echo "Referer: " . ($_SERVER['HTTP_REFERER'] ?? '') . "\n";
echo "Method: " . ($_SERVER['REQUEST_METHOD'] ?? '') . "\n";
echo "All headers:\n";
foreach (getallheaders() as $k => $v) {
    echo "$k: $v\n";
}

// Show whether the file exists and index.php content for sanity
echo "\nFiles present in project root:\n";
$files = scandir(__DIR__);
foreach ($files as $f) {
    echo $f . "\n";
}

?>
