import { initializeApp } from 'firebase/app'
import { getFirestore, connectFirestoreEmulator, collection } from 'firebase/firestore'

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
  connectFirestoreEmulator(db, 'localhost', 8080)
}
