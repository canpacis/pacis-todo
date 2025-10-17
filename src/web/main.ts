import Alpine from "alpinejs";
import { initializeApp } from "firebase/app";
import { getAuth, GoogleAuthProvider, signInWithPopup } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyA7uNtqoExbVEz7hEe_Sjcc5wU3e5VzXH8",
  authDomain: "pacis-todo.firebaseapp.com",
  projectId: "pacis-todo",
  storageBucket: "pacis-todo.firebasestorage.app",
  messagingSenderId: "99267217707",
  appId: "1:99267217707:web:3d861f8584506c1153dd6e",
};

const app = initializeApp(firebaseConfig);

const provider = new GoogleAuthProvider();
provider.addScope("profile");

const auth = getAuth(app);

Alpine.data("auth", () => ({
  async login() {
    try {
      const credentials = await signInWithPopup(auth, provider);
      const token = await credentials.user.getIdToken(true);
      window.location.replace(`/login/done?token=${token}`);
    } catch (error) {
      console.log("error", error);
    }
  },
  logout() {
    window.location.replace("/logout");
  },
}));

Alpine.start();
