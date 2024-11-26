"use client"

import * as React from "react"

import { cn, connectWallet } from "@/lib/utils"
import { API } from "@/lib/config"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Button } from "../ui/button"
import { LoaderIcon as Spinner, GithubIcon as GitHub } from "lucide-react"
import { toast, useToast } from "@/hooks/use-toast"

interface CreateAccountFormProps extends React.HTMLAttributes<HTMLDivElement> {}

export function CreateAccountForm({ className, ...props }: CreateAccountFormProps) {
  const [isLoading, setIsLoading] = React.useState<boolean>(false)

  async function onSubmit(event: React.SyntheticEvent) {
    event.preventDefault()
    setIsLoading(true)

    try {
      const username = event.target[0].value
      let res = await fetch(API.CHECK_USERNAME, {
        method: "POST",
        body: JSON.stringify({ Username: username }),
      })
      if (!res.ok) {
        toast({ title: "Something went wrong...", description: "Try again in a minute?" })
        return
      }

      let { exists } = await res.json()
      if (exists) {
        toast({ title: "Username already taken", description: "Try a different one?" })
        return
      }

      const handle = await connectWallet()
      if (handle) {
        await window.ethereum?.request({ method: "eth_requestAccounts" })

        const accounts = await handle.eth.getAccounts()
        const address = accounts[0]

        res = await fetch(API.GET_CHALLENGE, {
          method: "POST",
          body: JSON.stringify({ Username: username, Address: address }),
        })
        const body = await res.json()
        if (!res.ok) {
          console.error(body)
          throw new Error("Problem occured during registration")
        }

        // const signature = await signer.signMessage(body.challenge)
        console.log(body.challenge, address)
        // ERROR HERE
        const signature = await handle.eth.personal.sign(body.challenge, address, "")
        // ^^^^^^
        const payload = {
          Signature: signature,
          Username: username,
          Challenge: body.challenge,
        }

        res = await fetch(API.SIGN_UP, {
          method: "POST",
          body: JSON.stringify(payload),
        })

        console.log(res)
      } else {
        throw new Error("Unable to connect to your browser wallet")
      }
    } catch (e: any) {
      console.error(e)
      toast({
        title: "Something went wrong...",
        description: e.message || "Try again in a minute?",
      })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <form onSubmit={onSubmit}>
        <div className="grid gap-2">
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="email">
              Username
            </Label>
            <Input
              pattern="^[a-zA-z0-9]*$"
              title="(only alphanumeric characters allowed)"
              id="username"
              placeholder="PixelConnoisseur"
              type="text"
              autoCapitalize="none"
              autoCorrect="off"
              required
              disabled={isLoading}
            />
          </div>
          <Button disabled={isLoading}>
            {isLoading && <Spinner className="mr-2 h-4 w-4 animate-spin" />}
            Sign Up
          </Button>
        </div>
      </form>
    </div>
  )
}
