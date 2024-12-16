import React, { useState } from "react";

const CopilotUI = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");

  const handleSend = () => {
    if (input.trim() !== "") {
      setMessages([...messages, { text: input, sender: "user" }]);
      setInput("");
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
      <div className="bg-white shadow-md rounded-lg p-6 max-w-xl w-full">
        <h1 className="text-2xl font-bold mb-4">Data Analysis Chatbot</h1>
        <div className="mb-4 h-96 overflow-y-scroll">
          {messages.map((msg, index) => (
            <div
              key={index}
              className={`mb-2 text-${
                msg.sender === "user" ? "right" : "left"
              }`}
            >
              <div
                className={`inline-block p-2 rounded-lg ${
                  msg.sender === "user"
                    ? "bg-blue-500 text-white"
                    : "bg-gray-300 text-black"
                }`}
              >
                {msg.text}
              </div>
            </div>
          ))}
        </div>
        <div className="flex">
          <textarea
            className="w-full p-2 border border-gray-300 rounded-md mr-2"
            rows="2"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type your query here..."
          ></textarea>
          <button
            className="bg-blue-500 text-white p-2 rounded-md"
            onClick={handleSend}
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
};

export default CopilotUI;
