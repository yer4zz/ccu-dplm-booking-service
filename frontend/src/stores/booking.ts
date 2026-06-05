import { defineStore }  from 'pinia'
import { ref, computed } from 'vue'
import { servicesApi }  from '@/api/services'
import { mastersApi }   from '@/api/masters'
import { bookingsApi }  from '@/api/bookings'
import { loyaltyApi }   from '@/api/loyalty'
import type { Service, Master, Slot, BookingPriceCalc } from '@/types'

export type Step = 'service' | 'master' | 'slot' | 'confirm' | 'done'

export const useBookingStore = defineStore('booking', () => {
  const services         = ref<Service[]>([])
  const masters          = ref<Master[]>([])
  const slots            = ref<Slot[]>([])

  const selectedServices = ref<Service[]>([])
  const selectedMaster   = ref<Master | null>(null)
  const selectedSlot     = ref<Slot | null>(null)
  const selectedDate     = ref<string>('')
  const notes            = ref<string>('')
  const pointsToUse      = ref<number>(0)
  const priceCalc        = ref<BookingPriceCalc | null>(null)
  const vibeMode         = ref('')
  const vibeNote         = ref('')

  const step    = ref<Step>('service')
  const loading = ref(false)
  const error   = ref<string | null>(null)

  const mastersForServices = computed(() => {
    if (selectedServices.value.length === 0) return masters.value
    const ids = new Set(selectedServices.value.map(s => s.id))
    const filtered = masters.value.filter(m =>
      m.services.some(s => ids.has(s.id))
    )
    return filtered.length > 0 ? filtered : masters.value
  })

  const canProceed = computed(() => {
    if (step.value === 'service') return selectedServices.value.length > 0
    if (step.value === 'master')  return !!selectedMaster.value
    if (step.value === 'slot')    return !!selectedSlot.value
    return true
  })

  const primaryService = computed(() => selectedServices.value[0] ?? null)

  let cacheTimestamp = 0
  const CACHE_TTL    = 5 * 60 * 1000

  async function init() {
    if (Date.now() - cacheTimestamp < CACHE_TTL && services.value.length > 0) {
      return
    }

    loading.value = true
    try {
      const [s, m] = await Promise.all([
        servicesApi.list(),
        mastersApi.list(),
      ])
      services.value  = s ?? []
      masters.value   = m ?? []
      cacheTimestamp  = Date.now()
    } finally {
      loading.value = false
    }
  }

  function toggleService(svc: Service) {
    const idx = selectedServices.value.findIndex(s => s.id === svc.id)
    if (idx >= 0) {
      selectedServices.value.splice(idx, 1)
    } else {
      selectedServices.value.push(svc)
    }
    priceCalc.value = null
  }

  const slotsLoading = ref(false)

  async function loadSlots() {
    if (!selectedMaster.value || selectedServices.value.length === 0 || !selectedDate.value) return
    slotsLoading.value = true
    slots.value = []
    try {
      slots.value = await mastersApi.getSlots(
        selectedMaster.value.id,
        selectedServices.value.map(s => s.id),
        selectedDate.value,
      )
    } catch (e) {
      console.error('loadSlots error', e)
      slots.value = []
    } finally {
      slotsLoading.value = false
    }
  }

  let calcPriceTimer: ReturnType<typeof setTimeout> | null = null
  async function calcPrice() {
    if (!selectedMaster.value || selectedServices.value.length === 0) return

    if (calcPriceTimer) clearTimeout(calcPriceTimer)
    calcPriceTimer = setTimeout(async () => {
      try {
        priceCalc.value = await loyaltyApi.calcPrice({
          master_id:   selectedMaster.value!.id,
          service_ids: selectedServices.value.map(s => s.id),
          points_used: pointsToUse.value,
        })
      } catch (e) {
        console.error('calcPrice error', e)
      }
    }, 300)
  }

  async function confirm() {
    if (!selectedMaster.value || selectedServices.value.length === 0 || !selectedSlot.value) return

    loading.value = true
    error.value   = null

    const prevStep = step.value
    step.value = 'done'

    try {
      await bookingsApi.create({
        master_id:   selectedMaster.value.id,
        service_ids: selectedServices.value.map(s => s.id),
        starts_at:   selectedSlot.value.starts_at,
        notes:       notes.value,
        points_used: pointsToUse.value,
        vibe_mode:   vibeMode.value,
        vibe_note:   vibeNote.value,
      })
    } catch (e: any) {
      step.value  = prevStep
      error.value = e.response?.data?.error === 'slot_unavailable'
        ? 'Этот слот уже занят, выберите другое время'
        : 'Произошла ошибка, попробуйте ещё раз'
    } finally {
      loading.value = false
    }
  }

  function nextStep() {
    const order: Step[] = ['service', 'master', 'slot', 'confirm', 'done']
    const i = order.indexOf(step.value)
    if (i < order.length - 1) step.value = order[i + 1]
    if (step.value === 'confirm') calcPrice()
  }

  function prevStep() {
    const order: Step[] = ['service', 'master', 'slot', 'confirm', 'done']
    const i = order.indexOf(step.value)
    if (i > 0) step.value = order[i - 1]
  }

  function reset() {
    selectedServices.value = []
    selectedMaster.value   = null
    selectedSlot.value     = null
    selectedDate.value     = ''
    notes.value            = ''
    pointsToUse.value      = 0
    priceCalc.value        = null
    step.value             = 'service'
    error.value            = null
    vibeMode.value         = ''
    vibeNote.value         = ''
    cacheTimestamp         = 0
  }

  return {
    services, masters, slots,
    selectedServices, selectedMaster, selectedSlot, selectedDate,
    notes, pointsToUse, priceCalc, vibeMode, vibeNote,
    step, loading, error, slotsLoading,
    mastersForServices, canProceed, primaryService,
    init, toggleService, loadSlots, calcPrice, confirm,
    nextStep, prevStep, reset,
  }
})