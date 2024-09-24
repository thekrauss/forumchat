document.addEventListener("DOMContentLoaded", function () {
        // URLs pour la connexion WebSocket et API REST.
    const BASE_URL = "http://localhost:8079";
    const WEBSOCKET_URL = "ws://localhost:8079/ws";

    // Création d'une connexion WebSocket
    const ws = new WebSocket(WEBSOCKET_URL);

    const usersList = document.getElementById("users");
    const messageDiv = document.getElementById("messages");
    const messageForm = document.getElementById("message-form");
    const messageInput = document.getElementById("message-input");
    const typingIndicator = document.getElementById('typingIndicator'); 

    // État du chat pour stocker les données essentielles à la session
    const chatState = {
        messageHistory: {},
        currentChatUser: null,
        senderUsername: getUsernameFromCookie(),
        isLoading: false,
        typingTimeout: null,
        lastTypingSent: 0, 
    };
    console.log(typingIndicator); 


   typingIndicator.style.display = 'none';


    function getUsernameFromCookie() {
        const name = "username=";
        const decodedCookie = decodeURIComponent(document.cookie);
        const username = decodedCookie.split('; ').find(row => row.startsWith(name))?.split('=')[1] || "";
        console.log("Username from cookie:", username); 
        return username;
    }
    
    ws.onopen = () => console.log("Connected to the WebSocket server");

    ws.onerror = (error) => console.error("WebSocket error:", error);

    // Réception des messages via WebSocket.
    ws.onmessage = event => {
        const message = JSON.parse(event.data);
        console.log("Message received:", message);
    
        // Gestion des messages de type 'typing' pour l'indicateur de saisie.
        if (message.type === 'typing') {
            console.log("Handling typing event");
            handleTypingIndicator(message); 
        } 
            // Gestion des messages textuels.
        else if (message.type === 'message' && message.content && message.content.trim() !== "") {
            console.log("Handling new message");
            handleNewMessage(message);
        } else {
            console.log("Received an unknown event or empty message, ignored.");
        }
    };
    

    ws.onclose = () => console.log("Disconnected from the WebSocket server");


    // Fonction pour gérer l'affichage de l'indicateur de saisie.
    function handleTypingIndicator(message) {
        const typingIndicator = document.getElementById("typingIndicator");
    
        if (!typingIndicator) {
            console.error("L'élément typingIndicator n'a pas été trouvé !");
            return;
        }
    
        console.log("Typing event:", message);
    
        // Affichage de l'indicateur lorsque l'utilisateur commence à écrire.
        if (message.isTyping) {
            console.log(`${message.senderUsername} is typing...`);
            typingIndicator.textContent = `${message.senderUsername}`;

            typingIndicator.innerHTML += `
            <span></span><span></span><span></span>
            `;
    
            typingIndicator.style.display = 'block';
            
            typingIndicator.classList.add("show");
        } else {
            // Masquer l'indicateur lorsque l'utilisateur cesse de taper.
            console.log(`${message.senderUsername} stopped typing`);
            typingIndicator.classList.remove("show");
    
            setTimeout(() => {
                typingIndicator.style.display = 'none';
            }, 300); 
        }
    }
    
        // Fonction pour gérer la réception de nouveaux messages.
    function handleNewMessage(message) {
            // Ajout des messages à l'historique de l'expéditeur et du destinataire.
        if (!chatState.messageHistory[message.senderUsername]) {
            chatState.messageHistory[message.senderUsername] = [];
        }
        if (!chatState.messageHistory[message.targetUsername]) {
            chatState.messageHistory[message.targetUsername] = [];
        }

        chatState.messageHistory[message.senderUsername].push(message);
        chatState.messageHistory[message.targetUsername].push(message);
        
        // Affichage des messages uniquement si l'utilisateur sélectionné est impliqué dans la conversation.
        if (chatState.currentChatUser === message.senderUsername || chatState.currentChatUser === message.targetUsername) {
            displayMessages(chatState.currentChatUser);
        }
    }

    // Gestion de l'événement de saisie dans l'input de message pour envoyer un statut de saisie.
    messageInput.addEventListener('input', () => {
        clearTimeout(chatState.typingTimeout);
    
        const inputValue = messageInput.value.trim();
        if (inputValue === "") {
            sendTypingStatus(false);
            typingIndicator.style.display = 'none';
            return;
        }
    
        const timeSinceLastTyping = Date.now() - chatState.lastTypingSent;
        if (timeSinceLastTyping > 1000) {
            sendTypingStatus(true); // Indiquer que l'utilisateur est en train de taper.
            chatState.lastTypingSent = Date.now();
        }
    
                // Timeout pour envoyer un statut "ne tape plus" après 1.5 secondes d'inactivité.
        chatState.typingTimeout = setTimeout(() => {
            sendTypingStatus(false);
        }, 1500);
    });
    

    function sendTypingStatus(isTyping) {
        const message = {
            type: 'typing',
            senderUsername: chatState.senderUsername,
            targetUsername: messageInput.dataset.targetUsername || 'all',
            isTyping
        };
        console.log("Typing status sent:", message);
        ws.send(JSON.stringify(message));
    }

        // Gestion de la sélection d'un utilisateur pour démarrer un chat.
    usersList.addEventListener('click', event => {
        if (event.target.tagName === 'LI') {
            const targetUsername = event.target.textContent.trim();
            messageInput.dataset.targetUsername = targetUsername;
            chatState.currentChatUser = targetUsername;

            usersList.querySelectorAll('li').forEach(user => user.classList.remove('active'));
            event.target.classList.add('active');

            typingIndicator.textContent = ''; // Réinitialisation de l'indicateur de saisie.
            fetchMessageHistory(targetUsername); // Chargement de l'historique des messages.
        }
    });

        // Gestion de l'envoi du formulaire de message.
    messageForm.addEventListener("submit", event => {
        event.preventDefault();
        const messageContent = messageInput.value.trim();
        if (!messageContent) {
            console.log("Message content is empty");
            return;
        }
    
                // Création et envoi du message via WebSocket.
        const message = {
            type: 'message',  
            content: messageContent,
            targetUsername: messageInput.dataset.targetUsername || 'all',
            senderUsername: chatState.senderUsername,
            timestamp: new Date().toISOString()
        };
    
        console.log("Message sent:", message);
        ws.send(JSON.stringify(message));
    
        messageInput.value = ""; // Réinitialisation du champ input après envoi.
        typingIndicator.style.display = 'none';
        sendTypingStatus(false);  // Statut "ne tape plus" après envoi
    });
    

        // Fonction pour récupérer l'historique des messages d'un utilisateur.
    async function fetchMessageHistory(username, loadMore = false) {
        if (!loadMore && chatState.messageHistory[username]) {
            displayMessages(username, true);
            return;
        }
        
        if (chatState.isLoading) return;
        chatState.isLoading = true;
        const offset = loadMore ? chatState.messageHistory[username].length : 0;

        try {
            const response = await fetch(`${BASE_URL}/message?user=${username}&offset=${offset}`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            const messages = await response.json();
            if (!chatState.messageHistory[username]) chatState.messageHistory[username] = [];
            chatState.messageHistory[username] = loadMore ? [...messages, ...chatState.messageHistory[username]] : messages;
            displayMessages(username, !loadMore);
        } catch (error) {
            console.error("Error fetching message history:", error);
        } finally {
            chatState.isLoading = false;
        }
    }

        // Fonction pour afficher les messages de la conversation actuelle.
    function displayMessages(username, scrollToBottom = true) {
        const messages = chatState.messageHistory[username] || [];
        const messagesHTML = messages.map(message => {
            const time = new Date(message.timestamp).toLocaleTimeString();
            return `<div class="${message.senderUsername === chatState.senderUsername ? 'user-message' : 'other-message'}">
                        <span>[${time}] ${message.senderUsername}:</span> ${message.content}
                    </div>`;
        }).join("");
        
        messageDiv.innerHTML = messagesHTML;

            // Optionnel : défilement automatique vers le bas pour voir les nouveaux messages.
        if (scrollToBottom) {
            messageDiv.scrollTop = messageDiv.scrollHeight;
        }
    }

        // Chargement initial des utilisateurs en ligne.
    fetchOnlineUsers();

    async function fetchOnlineUsers() {
        try {
            const response = await fetch(`${BASE_URL}/user-list`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            const users = await response.json();
            updateUsersList(users);
        } catch (error) {
            console.error("Error fetching online users:", error);
            usersList.innerHTML = '<p class="error-message">Failed to load online users</p>';
        }
    }

        // Mise à jour de l'affichage de la liste des utilisateurs.
    function updateUsersList(users) {
        users.sort((a, b) => a.localeCompare(b));
        usersList.innerHTML = users.map(user => `<li>${user} <span class="online"></span></li>`).join("");
    }
});
