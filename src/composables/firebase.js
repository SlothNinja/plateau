import { initializeApp } from 'firebase/app'
import { getFirestore, connectFirestoreEmulator } from 'firebase/firestore'
import { getAuth, connectAuthEmulator } from 'firebase/auth'

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyDFiiIyKsp70-Im0CByK3zxGENGCxypG3w",
  authDomain: "plateau-fs-slothninja-games.firebaseapp.com",
  projectId: "plateau-fs-slothninja-games",
  storageBucket: "plateau-fs-slothninja-games.appspot.com",
  messagingSenderId: "467490981249",
  appId: "1:467490981249:web:ba1b9c621ad5a5d6f3d480"
};

// Initialize Firebase
export const firebaseApp = initializeApp(firebaseConfig);

// used for the firestore refs
export const db = getFirestore(firebaseApp)
if (process.env.NODE_ENV === 'development') {
  connectFirestoreEmulator(db, '127.0.0.1', 8080)
}

export const auth = getAuth(firebaseApp)
if (process.env.NODE_ENV === 'development') {
  connectAuthEmulator(auth, "http://127.0.0.1:9099")
}
