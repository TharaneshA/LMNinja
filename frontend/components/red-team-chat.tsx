"use client"

import { useState, useRef, useEffect } from 'react'
import { useAppStore } from '@/lib/state'
import { SendMessage } from '../../wailsjs/go/app/App'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { SendHorizonal, Bot, User, Loader2 } from 'lucide-react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { toast } from 'sonner'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

export function RedTeamChat() {
  const { activeModel } = useAppStore()
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const scrollAreaRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (scrollAreaRef.current) {
        const viewport = scrollAreaRef.current.querySelector('div[data-radix-scroll-area-viewport]')
        if(viewport) {
            viewport.scrollTop = viewport.scrollHeight;
        }
    }
  }, [messages, isLoading]);

  const handleSendMessage = async () => {
    if (!input.trim() || isLoading) return

    const userMessage: Message = { role: 'user', content: input }
    setMessages((prev) => [...prev, userMessage])
    setInput('')
    setIsLoading(true)

    try {
      const response = await SendMessage(input)
      const botMessage: Message = { role: 'assistant', content: response }
      setMessages((prev) => [...prev, botMessage])
    } catch (error) {
      toast.error('Failed to get response', { description: String(error) })
      setMessages((prev) => prev.slice(0, -1))
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Card className="flex-1 flex flex-col bg-transparent border-none shadow-none mt-4 h-full">
      <CardContent className="flex-1 p-0 overflow-hidden">
        <ScrollArea className="h-full p-4" ref={scrollAreaRef}>
          <div className="space-y-6">
            {messages.map((msg, index) => (
              <div key={index} className={`flex items-start gap-4 ${msg.role === 'user' ? 'justify-end' : ''}`}>
                {msg.role === 'assistant' && (<div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center flex-shrink-0"><Bot className="w-5 h-5 text-primary-foreground" /></div>)}
                <div className={`max-w-[75%] p-3 rounded-lg ${msg.role === 'user' ? 'bg-[#4A90E2] text-white' : 'bg-[#1B2636]'}`}><p className="text-sm whitespace-pre-wrap">{msg.content}</p></div>
                {msg.role === 'user' && (<div className="w-8 h-8 rounded-full bg-gray-600 flex items-center justify-center flex-shrink-0"><User className="w-5 h-5 text-white" /></div>)}
              </div>
            ))}
            {isLoading && (<div className="flex items-start gap-4"><div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center flex-shrink-0"><Bot className="w-5 h-5 text-primary-foreground" /></div><div className="max-w-[75%] p-3 rounded-lg bg-[#1B2636] flex items-center"><Loader2 className="w-5 h-5 animate-spin text-gray-400" /></div></div>)}
          </div>
        </ScrollArea>
      </CardContent>
      <CardFooter className="p-4 border-t border-[#2A3B52]">
        <div className="flex w-full items-center space-x-2">
          <Textarea placeholder={`Chat with ${activeModel?.Name}... (Shift+Enter for newline)`} className="flex-1 bg-[#0F1419] border-[#2A3B52] resize-none" value={input} onChange={(e) => setInput(e.target.value)} onKeyDown={(e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSendMessage(); }}}/>
          <Button onClick={handleSendMessage} disabled={isLoading || !input.trim()} className="bg-[#4A90E2] hover:bg-[#3A7BC8]"><SendHorizonal className="w-4 h-4" /></Button>
        </div>
      </CardFooter>
    </Card>
  )
}