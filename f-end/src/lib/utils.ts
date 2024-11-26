import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
// import { ethers, type BrowserProvider } from "ethers"
import Web3 from "web3"
import { type RegisteredSubscription } from "web3-eth"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export async function connectWallet(): Promise<Web3<RegisteredSubscription> | void> {
  if (window.ethereum == null) {
    throw new Error("Wallet isn't installed")
  }

  // return new ethers.BrowserProvider(window.ethereum)
  return new Web3(window.ethereum)
}
