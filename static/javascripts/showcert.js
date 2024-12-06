console.log("showcert.js loaded"); // Debugging: Ensure the script is loaded

async function fetchCertificate(event) {
    event.preventDefault(); // Prevent the default form submission behavior
    console.log("fetchCertificate called"); // Debugging: Ensure the function is executed

    const domain = document.getElementById("domain").value;
    const resultDiv = document.getElementById("result");

    // Clear previous results
    resultDiv.innerHTML = "Fetching data...";

    try {
        // Fetch data from the server
        const response = await fetch(`/fetch-cert?domain=${encodeURIComponent(domain)}`);
        console.log("Response received:", response); // Debugging: Log the response

        const data = await response.text();


        if (response.ok) {
            // Modify the data to make "Subject" and "Issuer" bold
            const formattedData = data
                .replace(/(Subject:.*)/g, "<b>$1</b>") // Bold the Subject line
                .replace(/(Issuer:.*)/g, "<b>$1</b>") // Bold the Issuer line
                .replace(/(Failed to.*)/g, "<b class=\"bold-red\">$1</b>");
            // Display the result in the result div
            resultDiv.innerHTML = `<pre>${formattedData}</pre>`;
        } else {
            resultDiv.innerHTML = `<span style="color: red;">Error: ${data}</span>`;
        }

            
    } catch (error) {
        console.error("Error fetching certificate:", error); // Debugging: Log any error
        resultDiv.innerHTML = `<span style="color: red;">Failed to fetch data: ${error.message}</span>`;
    }
}

