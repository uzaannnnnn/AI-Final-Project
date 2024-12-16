import React, { useState, useEffect } from "react";
import axios from "axios";
import "./index.css";

function App() {
  const [file, setFile] = useState(null);
  const [fileName, setFileName] = useState("");
  const [query, setQuery] = useState("");
  const [messages, setMessages] = useState([]);
  const [response, setResponse] = useState("");
  const [displayedResponse, setDisplayedResponse] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (response && !loading) {
      console.log("Response received:", response);
      const words = response.split(" ");
      let index = 0;
      let currentResponse = "";

      const intervalId = setInterval(() => {
        if (index < words.length) {
          currentResponse += (index > 0 ? " " : "") + words[index];
          setDisplayedResponse(currentResponse);
          index++;
        } else {
          clearInterval(intervalId);
        }
      }, 100); // Ubah interval sesuai kebutuhan

      return () => clearInterval(intervalId);
    }
  }, [response, loading]);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
    setFileName(selectedFile ? selectedFile.name : "");
  };

  const handleSend = async () => {
    if (query.trim() === "") return;

    setLoading(true); // Mulai loading

    setMessages((prevMessages) => [
      ...prevMessages,
      { text: query, sender: "user" },
      { text: "", sender: "bot", loading: true },
    ]);

    let res;
    try {
      if (file) {
        const formData = new FormData();
        formData.append("file", file);
        formData.append("query", query);

        res = await axios.post("http://localhost:8000/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });
        console.log("File Upload Response:", res.data); // Debugging
        setFile(null); // Reset file after upload
        setFileName(""); // Reset file name after upload
      } else {
        res = await axios.post("http://localhost:8000/chat", { query });
        console.log("Chat Response:", res.data); // Debugging
      }

      const serverResponse = res.data?.answer || "No response from server";
      console.log("Setting response:", serverResponse);

      setMessages((prevMessages) =>
        prevMessages.map((msg, index) => {
          if (index === prevMessages.length - 1 && msg.sender === "bot") {
            return { text: serverResponse, sender: "bot", loading: false };
          }
          return msg;
        })
      );
      setResponse(serverResponse);
      setQuery("");
      setDisplayedResponse(""); // Reset displayedResponse before starting
    } catch (error) {
      console.error("Error:", error);
    } finally {
      setLoading(false); // Selesai loading
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter") {
      handleSend();
    }
  };

  const formatResponse = (response) => {
    if (typeof response !== "string") {
      response = JSON.stringify(response, null, 2);
    }
    return response.split("\n").map((line, index) => (
      <p key={index} className="mb-2">
        {line}
      </p>
    ));
  };

  return (
    <div className="min-h-screen bg-gray-800 flex flex-col items-center justify-center p-4">
      <div className="max-w-[1000px] w-full">
        <h1 className="text-2xl font-bold mb-4 text-white uppercase text-center">
          Data Analysis Chatbot
        </h1>
        <div className="mb-4 p-4 pb-12 rounded-tl-3xl rounded-tr-3xl border-l border-t border-r border-gray-400 h-[646px] overflow-y-auto">
          {messages.map((msg, index) => (
            <div
              key={index}
              className={`mb-2 flex ${
                msg.sender === "user" ? "justify-end" : "justify-start"
              }`}
            >
              <div
                className={`inline-block px-2 pt-1 rounded-lg ${
                  msg.sender === "user"
                    ? "bg-gray-700 text-white font-sans font-semibold"
                    : "text-white font-sans font-semibold "
                }`}
              >
                {msg.sender === "bot" && msg.loading ? (
                  <div className="flex items-center">
                    <span className="mr-2">Loading</span>
                    <div className="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-blue-500"></div>
                  </div>
                ) : msg.sender === "bot" && index === messages.length - 1 ? (
                  displayedResponse
                ) : (
                  formatResponse(msg.text)
                )}
              </div>
            </div>
          ))}
        </div>
        <div className="w-3/4 mx-auto">
          <div className="flex items-center border border-gray-700 py-2 px-5 rounded-full w-1/2 mx-auto fixed bottom-0 mb-4 ">
            <input
              type="file"
              onChange={handleFileChange}
              className="hidden"
              id="file-upload"
            />
            <label htmlFor="file-upload" className="cursor-pointer mr-4">
              <svg
                width="29px"
                height="29px"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                stroke="#000000"
              >
                <g id="SVGRepo_bgCarrier" strokeWidth="0"></g>
                <g
                  id="SVGRepo_tracerCarrier"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                ></g>
                <g id="SVGRepo_iconCarrier">
                  {" "}
                  <path
                    d="M20 14V7C20 5.34315 18.6569 4 17 4H12M20 14L13.5 20M20 14H15.5C14.3954 14 13.5 14.8954 13.5 16V20M13.5 20H7C5.34315 20 4 18.6569  4 17V12"
                    stroke="#ffffff"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  ></path>{" "}
                  <path
                    d="M7 4V7M7 10V7M7 7H4M7 7H10"
                    stroke="#ffffff"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  ></path>{" "}
                </g>
              </svg>
            </label>
            {fileName && (
              <span className="w-28 pl-1 text-xs font-sans font-semibold mr-2 rounded-full bg-white">
                {fileName}
              </span>
            )}
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder={
                file ? "Describe your upload..." : "Ask a question..."
              }
              className="py-2 px-5 rounded-full bg-gray-900 w-full placeholder:text-lg text-white text-lg"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
