import { runPythonCode } from "https://esm.town/v/iamseeley/pyodideMod";

document.getElementById('run-python').addEventListener('click', async () => {
    const pythonCode = `
def greet(name):
    return f"Hello, {name}!"

result = greet("JavaScript")
result
    `;

    try {
        const result = await runPythonCode(pythonCode);
        document.getElementById('output').innerText = `Python Output: ${result}`;
    } catch (error) {
        console.error("Error executing Python code:", error);
        document.getElementById('output').innerText = `Error: ${error.message}`;
    }
});
