// Function to fetch and display certificate information
async function fetchCertificate(event) {
    event.preventDefault(); // Prevent form submission

    const domain = document.getElementById("domain").value;
    const resultDiv = document.getElementById("result");

    // Clear previous results
    resultDiv.innerHTML = "Fetching data...";

    try {
        // Fetch data from the server
        const response = await fetch(`/fetch-cert?domain=${encodeURIComponent(domain)}`);
        const data = await response.text();

        if (response.ok) {
            // Display the result in the result div
            resultDiv.innerHTML = `<pre>${data}</pre>`;
        } else {
            resultDiv.innerHTML = `<span style="color: red;">Error: ${data}</span>`;
        }
    } catch (error) {
        resultDiv.innerHTML = `<span style="color: red;">Failed to fetch data: ${error}</span>`;
    }
}