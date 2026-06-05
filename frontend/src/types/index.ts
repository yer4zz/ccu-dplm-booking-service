export interface Service {
  id: string
  name: string
  name_kz?: string
  name_en?: string
  description?: string
  desc_kz?: string
  desc_en?: string
  category: string
  duration_min: number
  price: number
  is_active: boolean
}
  
  export interface Master {
    id: string
    full_name: string
    bio: string
    experience_years: number
    avatar_url: string
    instagram: string
    services: Service[]
    is_active: boolean
  }
  
  export interface Slot {
    starts_at: string 
    ends_at: string
  }
  
  export interface Booking {
    id: string
    master_id: string
    service_id: string
    client_id: string
    starts_at: string
    ends_at: string
    status: 'pending' | 'confirmed' | 'cancelled' | 'completed' | 'no_show'
    price_paid: number
    notes: string
    created_at: string
    service_name?: string
    master_name?: string
    client_name?: string 
  }
  
  export interface User {
    id: string
    email: string
    full_name: string
    phone:     string
    role: 'client' | 'master' | 'admin'
  }

export interface LoyaltyAccount {
  client_id: string
  balance: number
  total_earned: number
  updated_at: string
}

export interface PointTransaction {
  id: string
  type: 'earn' | 'redeem' | 'expire' | 'refund'
  amount: number
  description: string
  created_at: string
}

export interface BookingPriceCalc {
  base_price: number
  service_count: number
  is_first_booking: boolean
  first_discount: number
  multi_discount: number
  points_available: number
  points_used: number
  points_discount: number
  total_discount: number
  final_price: number
}

export interface WaitlistEntry {
  id: string
  master_id: string
  service_id: string
  status: string
  created_at: string
}

export interface RescheduleOffer {
  id: string
  original_booking_id: string
  status: string
  offered_at: string
  expires_at: string
}


export interface LoyaltyAccount {
  client_id: string
  balance: number
  total_earned: number
  updated_at: string
}

export interface PointTransaction {
  id: string
  type: 'earn' | 'redeem' | 'expire' | 'refund'
  amount: number
  description: string
  created_at: string
}

export interface BookingPriceCalc {
  base_price: number
  service_count: number
  is_first_booking: boolean
  first_discount: number
  multi_discount: number
  points_available: number
  points_used: number
  points_discount: number
  total_discount: number
  final_price: number
}

export interface WaitlistEntry {
  id: string
  master_id: string
  service_id: string
  status: string
  created_at: string
}

export interface RescheduleOffer {
  id: string
  original_booking_id: string
  status: string
  offered_at: string
  expires_at: string
}


export interface BeautyStreak {
  current_streak:      number
  longest_streak:      number
  last_visit_date:     string | null
  streak_deadline:     string | null
  level:               'bronze' | 'silver' | 'gold' | 'platinum'
  days_until_deadline: number
  next_level_at:       number
  progress:            number
  streak_discount:     number
}

export interface SOSRequest {
  id:                    string
  master_id:             string
  service_id:            string
  preferred_range_start: string
  preferred_range_end:   string
  base_price:            number
  sos_price:             number
  status:                'pending' | 'accepted' | 'declined' | 'expired'
  client_note:           string
  master_note:           string
  expires_at:            string
  created_at:            string
  client_name?:          string
  service_name?:         string
}